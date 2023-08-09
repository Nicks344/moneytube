import mobx from 'mobx'
import { observable, autorun, action, decorate } from "mobx";
import { getBugReports } from 'api';
import { deleteReport } from 'api';
import { resolveReport } from 'api';

export default class BugReportsStore {
    reports = []

    constructor() {
        autorun(() => console.log(this.report));

        getBugReports().then(reports => this.reports = reports).catch(err => window.notify(`Ошибка получения баг-репортов: ${err}`))
    }

    deleteReport = id => {
        if (window.confirm(`Вы действительно хотите удалить баг-репорт ${id}?`)) {
            deleteReport(id).then(() => this.reports = this.reports.filter(u => u.ID !== id)).catch(err => window.notify(`Ошибка удаления баг-репорта: ${err}`))
        }
    }
    
    resolveReport = id => {
        if (window.confirm(`Вы действительно хотите пометить баг-репорт ${id} выполненым?`)) {
            resolveReport(id).then(() => {
                this.reports.find(u => u.ID !== id).Resolved = true
            }).catch(err => window.notify(`Ошибка: ${err}`))
        }
    }
}

decorate(BugReportsStore, {
    reports: observable,
});