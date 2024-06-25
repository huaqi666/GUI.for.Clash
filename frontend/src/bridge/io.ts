import { App } from '@wails/guiforcores/bridge'

type IOOptions = {
  Mode?: 'Binary' | 'Text'
}

const enum BridgeHTTPApi {
  WriteFile = '/bridge/fs/write',
  ReadFile = '/bridge/fs/read'
}

export const Writefile = async (path: string, content: string, options: IOOptions = {}) => {
  // const { flag, data } = await App.Writefile(path, content, { Mode: 'Text', ...options })
  const res = await fetch(BridgeHTTPApi.WriteFile, {
    method: 'POST',
    body: JSON.stringify({
      Path: path,
      Content: content,
      Options: { Mode: 'Text', ...options }
    })
  })
  const { flag, data } = await res.json()
  if (!flag) {
    throw data
  }
  return data
}

export const Readfile = async (path: string, options: IOOptions = {}) => {
  // const { flag, data } = await App.Readfile(path, { Mode: 'Text', ...options })
  const res = await fetch(BridgeHTTPApi.ReadFile, {
    method: 'POST',
    body: JSON.stringify({
      Path: path,
      Content: '',
      Options: { Mode: 'Text', ...options }
    })
  })
  const { flag, data } = await res.json()
  if (!flag) {
    throw data
  }
  return data
}

export const Movefile = async (source: string, target: string) => {
  const { flag, data } = await App.Movefile(source, target)
  if (!flag) {
    throw data
  }
  return data
}

export const Removefile = async (path: string) => {
  const { flag, data } = await App.Removefile(path)
  if (!flag) {
    throw data
  }
  return data
}

export const Copyfile = async (source: string, target: string) => {
  const { flag, data } = await App.Copyfile(source, target)
  if (!flag) {
    throw data
  }
  return data
}

export const FileExists = async (path: string) => {
  const { flag, data } = await App.FileExists(path)
  if (!flag) {
    throw data
  }
  return data === 'true'
}

export const AbsolutePath = async (path: string) => {
  const { flag, data } = await App.AbsolutePath(path)
  if (!flag) {
    throw data
  }
  return data
}

export const Makedir = async (path: string) => {
  const { flag, data } = await App.Makedir(path)
  if (!flag) {
    throw data
  }
  return data
}

export const Readdir = async (path: string) => {
  const { flag, data } = await App.Readdir(path)
  if (!flag) {
    throw data
  }
  return data
    .split('|')
    .filter((v) => v)
    .map((v) => {
      const [name, size, isDir] = v.split(',')
      return { name, size: Number(size), isDir: isDir === 'true' }
    })
}

export const UnzipZIPFile = async (path: string, output: string) => {
  const { flag, data } = await App.UnzipZIPFile(path, output)
  if (!flag) {
    throw data
  }
  return data
}

export const UnzipGZFile = async (path: string, output: string) => {
  const { flag, data } = await App.UnzipGZFile(path, output)
  if (!flag) {
    throw data
  }
  return data
}
