import request from '../utils/request'

const api = {
  admin: '/admin',
  res: '/res',
  def: '/defs',
}

export default api

export function listDef () {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'listDef'}
  })
}
export function getDef (id) {
  const data = {'action': 'getDef', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function saveDef (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDef', 'data': data}
  })
}
export function saveDefDesign (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDefDesign', 'data': data}
  })
}
export function removeDef (id) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'removeDef', id: id}
  })
}

export function getDefFieldTree (id) {
  const data = {'action': 'getDefFieldTree', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getDefField (id) {
  const data = {'action': 'getDefField', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function createDefField (targetId, mode) {
  const data = {'action': 'createDefField', id: targetId, mode: mode}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function saveDefField (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDefField', 'data': data}
  })
}

export function removeDefField (id) {
  const data = {'action': 'removeDefField', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function moveDefField (src, dist, mode) {
  const data = {'action': 'moveDefField', src: src, dist: dist, mode: ''+mode}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listDefFieldSection (fieldId) {
  const data = {'action': 'listDefFieldSection', id: fieldId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function createDefFieldSection (fieldId, sectionId) {
  const data = {'action': 'createDefFieldSection', data: { fieldId: ''+fieldId, sectionId: ''+sectionId}}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function updateDefFieldSection (section) {
  const data = {'action': 'updateDefFieldSection', data: section}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function removeDefFieldSection (sectionId) {
  const data = {'action': 'removeDefFieldSection', id: sectionId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function getDefFieldRefer (fieldId, resType) {
  const data = {'action': 'getDefFieldRefer', id: fieldId, mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function updateDefFieldRefer (refer) {
  const data = {'action': 'updateDefFieldRefer', data: refer}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function listDefFieldReferType (resType) {
  const data = {'action': 'listDefFieldReferType', mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function listDefFieldReferField (refer) {
  const data = {'action': 'listDefFieldReferField', data: refer}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listRanges () {
  const data = {'action': 'listRanges'}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getRanges (id) {
  const data = {'action': 'getRanges', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveRanges (model) {
  const data = {'action': 'saveRanges', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeRanges (id) {
  const data = {'action': 'removeRanges', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getResRangesItemTree (id) {
  const data = {'action': 'getResRangesItemTree', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getResRangesItem (id) {
  const data = {'action': 'getResRangesItem', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function createResRangesItem (rangesId, mode) {
  const data = {'action': 'createResRangesItem', domainId: rangesId, mode: mode}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveRangesItem (model) {
  const data = {'action': 'saveRangesItem', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}export function removeResRangesItem (itemId, rangesId) {
  const data = {'action': 'removeResRangesItem', id: itemId, domainId: rangesId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listText () {
  const data = {'action': 'listText'}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getText (id) {
  const data = {'action': 'getText', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveText (model) {
  const data = {'action': 'saveText', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeText (id) {
  const data = {'action': 'removeText', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listConfig () {
  const data = {'action': 'listConfig'}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getConfig (id) {
  const data = {'action': 'getConfig', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveConfig (model) {
  const data = {'action': 'saveConfig', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeConfig (id) {
  const data = {'action': 'removeConfig', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
