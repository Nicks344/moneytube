import Button from '@material-ui/core/Button';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import CreateIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import { observer } from 'mobx-react';
import * as moment from 'moment';
import React from 'react';
import { licenseStore } from 'stores';
import CreateUserDialog from './CreateUserDialog';

class License extends React.Component {
    render() {
        const { createNewUser, users, deleteUser } = licenseStore
        return (
            <>
                <CreateUserDialog />
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Ник</TableCell>
                            <TableCell>Версия</TableCell>
                            <TableCell>Ключ</TableCell>
                            <TableCell>Создан</TableCell>
                            <TableCell>Активирован</TableCell>
                            <TableCell>Кол-во дней</TableCell>
                            <TableCell>Осталось дней</TableCell>
                            <TableCell>HWID</TableCell>
                            <TableCell><Button color="primary" onClick={createNewUser}><CreateIcon /></Button></TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {users.sort((a, b) => a.createdAt < b.createdAt ? 1 : -1).map((user, i) =>
                            <TableRow style={{ backgroundColor: user.isActivated && !user.isActive ? 'rgba(255,0,0, 0.1)' : '' }} key={i}>
                                <TableCell>{user.name}</TableCell>
                                <TableCell>{user.version ? user.version : '1'}</TableCell>
                                <TableCell>{user.key}</TableCell>
                                <TableCell>{moment(user.createdAt).format('YYYY-MM-DD')}</TableCell>
                                <TableCell>{!user.isActivated ? '-' : moment(user.activatedAt).format('YYYY-MM-DD')}</TableCell>
                                <TableCell>{user.days}</TableCell>
                                <TableCell>{!user.isActivated || !user.isActive ? '-' : user.daysLeft}</TableCell>
                                <TableCell>{user.hwid}</TableCell>
                                <TableCell><Button color="primary" onClick={() => deleteUser(user)}><DeleteIcon style={{ fontSize: '20px' }} /></Button></TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </>
        )
    }
}

export default observer(License)