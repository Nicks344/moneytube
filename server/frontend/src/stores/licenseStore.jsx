import { createUser, deleteUser, getUsers } from 'api';
import { saveAs } from 'file-saver';
import { action, autorun, decorate, observable } from "mobx";

const initUserObj = Object.freeze({
    name: '',
    days: 0,
    daysReactivate: 0,
    version: '',
    count: 1
})

export default class LicenseStore {
    users = []
    createUserOpened = false
    newUser = null

    constructor() {
        autorun(() => console.log(this.report));

        getUsers().then(users => this.users = users).catch(err => window.notify(`Ошибка получения пользователей: ${err}`))
    }

    createNewUser = () => {
        this.createUserOpened = true
        this.newUser = Object.assign({}, initUserObj)
    }

    editUserField = (field, value) => {
        this.newUser[field] = value
    }

    submitNewUser = () => {
        if (!this.validate()) {
            return
        }

        createUser(this.newUser).then(users => {
            this.users.push(...users)
            this.createUserOpened = false
            var blob = new Blob([users.map(({ key }) => key).join('\n')], {
                type: "text/plain;charset=utf-8;",
            })
            saveAs(blob, "keys.txt")
        }).catch(err => window.notify(`Ошибка создания пользователя: ${err}`))
    }

    cancelNewUser = () => {
        this.createUserOpened = false
    }

    deleteUser = user => {
        if (window.confirm(`Вы действительно хотите удалить пользователя ${user.name}?`)) {
            deleteUser(user.key).then(() => this.users = this.users.filter(u => u.key !== user.key)).catch(err => window.notify(`Ошибка удаления пользователя: ${err}`))
        }
    }

    validate = () => {
        if (!this.newUser.name) {
            window.notify('Укажите имя пользователя', 'error')
            return false
        }

        if (this.newUser.days <= 0) {
            window.notify('Укажите кол-во дней больше нуля', 'error')
            return false
        }

        if (this.newUser.daysReactivate <= 0) {
            window.notify('Укажите кол-во дней для реактивации больше нуля', 'error')
            return false
        }

        return true
    }
}

decorate(LicenseStore, {
    users: observable,
    createUserOpened: observable,
    createNewUser: action,
    editUserField: action,
    submitNewUser: action
});