import React, { useState } from 'react'
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Button from '@material-ui/core/Button'
import { bugReportsStore } from 'stores'
import { observer } from 'mobx-react'
import * as moment from 'moment'
import ResolveIcon from '@material-ui/icons/Done'
import DeleteIcon from '@material-ui/icons/Delete'
import DownloadIcon from '@material-ui/icons/Archive'
import ButtonGroup from '@material-ui/core/ButtonGroup';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import { downloadReportData } from 'api';

export default observer(function BugReports() {
    const [desc, setDesc] = useState({
        opened: false,
        body: ''
    })

    const handleCloseDesc = () => {
        setDesc({
            opened: false,
            body: ''
        })
    }

    const handleOpenDesc = desc => {
        setDesc({
            opened: true,
            body: desc
        })
    }

    const { reports, deleteReport, resolveReport } = bugReportsStore

    return (
        <>
            <Dialog
                open={desc.opened}
                onClose={handleCloseDesc}
            >
                <DialogTitle id="simple-dialog-title">Описание ошибки</DialogTitle>
                <DialogContent>
                    <DialogContentText style={{ whiteSpace: 'pre-line' }}>{desc.body}</DialogContentText>
                </DialogContent>
            </Dialog>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Создан</TableCell>
                        <TableCell>Ник</TableCell>
                        <TableCell>Ошибка</TableCell>
                        <TableCell>Описание</TableCell>
                        <TableCell></TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {reports.map((report, i) =>
                        <TableRow style={{backgroundColor: report.Resolved ? 'rgba(0,255,0, 0.1)' : ''}} key={i}>
                            <TableCell>{report.ID}</TableCell>
                            <TableCell>{moment(report.CreatedAt).format('YYYY-MM-DD HH:mm')}</TableCell>
                            <TableCell>{report.UserName}</TableCell>
                            <TableCell>{report.Error}</TableCell>
                            <TableCell><Button color="primary" onClick={() => handleOpenDesc(report.Description)}>Просмотр</Button></TableCell>
                            <TableCell>
                                <ButtonGroup color="primary">
                                    <Button color="primary" onClick={() => resolveReport(report.ID)}><ResolveIcon style={{ fontSize: '20px' }} /></Button>
                                    <Button color="primary" onClick={() => downloadReportData(report.ID)}><DownloadIcon style={{ fontSize: '20px' }} /></Button>
                                    <Button color="primary" onClick={() => deleteReport(report.ID)}><DeleteIcon style={{ fontSize: '20px' }} /></Button>
                                </ButtonGroup>
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </>
    )
})