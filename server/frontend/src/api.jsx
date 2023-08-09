const host = window.location.origin.replace('5000', '4444')

export async function getUsers() {
	return await request('GET', `${host}/api/admin/users`)
}

export async function createUser(user) {
	return await request('POST', `${host}/api/admin/user`, user)
}

export async function deleteUser(key) {
	return await request('DELETE', `${host}/api/admin/user/${key}`)
}

export async function getBugReports() {
	return await request('GET', `${host}/api/admin/bugreports`)
}

export async function downloadReportData(id) {
	var win = window.open(`${host}/api/admin/bugreports/${id}`, '_blank')
	win.focus()
}

export async function resolveReport(id) {
	return await request('POST', `${host}/api/admin/bugreports/${id}/resolve`)
}

export async function deleteReport(id) {
	return await request('DELETE', `${host}/api/admin/bugreports/${id}`)
}

async function request(method, path, data) {
	let params = {
		method: method
	}
	if (data) {
		params.body = JSON.stringify(data)
		params.headers = {
			'Content-Type': 'application/json'
		}
	}

	const resp = await fetch(path, params)
	const ans = await resp.json()
	if (resp.status !== 200) {
		throw new Error(ans.error)
	}
	return ans.result
}