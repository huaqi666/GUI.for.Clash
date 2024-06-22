import { App } from '@wails/guiforcores/bridge'

import { Events } from '@wailsio/runtime'

type RequestType = {
  id: string
  method: string
  url: string
  headers: Record<string, string>
  body: string
}

type ResponseType = {
  status: number
  headers: Record<string, string>
  body: string
  options: { Mode: 'Binary' | 'Text' }
}

type HttpServerHandler = (
  req: RequestType,
  res: {
    end: (
      status: ResponseType['status'],
      headers: ResponseType['headers'],
      body: ResponseType['body'],
      options: ResponseType['options']
    ) => void
  }
) => Promise<void>

export const StartServer = async (address: string, id: string, handler: HttpServerHandler) => {
  const { flag, data } = await App.StartServer(address, id)
  if (!flag) {
    throw data
  }

  Events.On(id, async ({ data }: WailsEventsResponse<RequestType>) => {
    const { id, method, url, headers, body } = data
    try {
      await handler(
        {
          id,
          method,
          url,
          headers: Object.entries(headers).reduce((p, c: any) => ({ ...p, [c[0]]: c[1][0] }), {}),
          body
        },
        {
          end: (status, headers, body, options = { Mode: 'Text' }) => {
            Events.Emit({
              name: id,
              data: { status, headers, body, options }
            })
          }
        }
      )
    } catch (err: any) {
      console.log('Server handler err:', err, id)
      Events.Emit({
        name: id,
        data: {
          status: 500,
          headers: { 'Content-Type': 'text/plain; charset=utf-8' },
          body: err.message || err,
          options: { Mode: 'Text' }
        }
      })
    }
  })
  return { close: () => StopServer(id) }
}

export const StopServer = async (serverID: string) => {
  const { flag, data } = await App.StopServer(serverID)
  if (!flag) {
    throw data
  }
  Events.Off(serverID)
  return data
}

export const ListServer = async () => {
  const { flag, data } = await App.ListServer()
  if (!flag) {
    throw data
  }
  return data.split('|').filter((id) => id.length)
}
