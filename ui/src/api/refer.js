import request from '../utils/request'
import api from './manage'

export function getRefer (ownerId, ownerType) {
  const data = {'action': 'getRefer', id: ownerId, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function updateRefer (refer, ownerType) {
  const data = {'action': 'updateRefer', data: refer, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

// selection input on page
export function listReferFileForSelection (resType) {
  const data = {'action': 'listReferFileForSelection', mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function listReferSheetForSelection (resType) {
  const data = {'action': 'listReferSheetForSelection', mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function listReferFieldForSelection (referId, referType) {
  const data = {'action': 'listReferFieldForSelection', id: referId, mode: referType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
