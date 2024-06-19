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
  end: (
    status: number,
    headers: Record<string, string>,
    body: string,
    options: { mode: 'Binary' | 'Text' }
  ) => void
}
type HttpServerHandler = (req: RequestType, res: ResponseType) => Promise<void>

export const StartServer = async (address: string, id: string, handler: HttpServerHandler) => {
  const { flag, data } = await App.StartServer(address, id)
  if (!flag) {
    throw data
  }

  Events.On(id, async (...args) => {
    const [id, method, url, headers, body] = args
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
          end: (status, headers, body, options = { mode: 'Text' }) => {
            Events.Emit(id, status, JSON.stringify(headers), body, JSON.stringify(options))
          }
        }
      )
    } catch (err: any) {
      console.log('Server handler err:', err, id)
      Events.Emit(
        id,
        500,
        JSON.stringify({ 'Content-Type': 'text/plain; charset=utf-8' }),
        err.message || err,
        JSON.stringify({ Mode: 'Text' })
      )
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
