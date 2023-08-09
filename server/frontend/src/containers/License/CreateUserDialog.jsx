import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import TextField from '@material-ui/core/TextField';
import { observer } from 'mobx-react';
import React from 'react';
import { licenseStore } from 'stores';

class CreateUserDialog extends React.Component {
    render() {
        const { createUserOpened, submitNewUser, editUserField, cancelNewUser } = licenseStore
        return (
            <Dialog open={createUserOpened} onClose={cancelNewUser}>
                <DialogTitle>Добавить</DialogTitle>
                <DialogContent>
                    <TextField
                        autoFocus
                        label="Имя пользователя"
                        onChange={e => editUserField('name', e.target.value)}
                        fullWidth
                    />
                    <TextField
                        style={{ marginTop: 10 }}
                        label="Кол-во дней"
                        type="number"
                        onChange={e => editUserField('days', parseInt(e.target.value))}
                        InputProps={{
                            inputProps: {
                                min: 1
                            }
                        }}
                        fullWidth
                    />
                    <TextField
                        style={{ marginTop: 10 }}
                        label="Кол-во дней для реактивации"
                        type="number"
                        onChange={e => editUserField('daysReactivate', parseInt(e.target.value))}
                        InputProps={{
                            inputProps: {
                                min: 1
                            }
                        }}
                        fullWidth
                    />

                    <TextField
                        style={{ marginTop: 10 }}
                        label="Версия(только 2 или пусто)"
                        onChange={e => e.target.value == '' || e.target.value == '2' ? editUserField('version', e.target.value) : null}
                        fullWidth
                    />

                    <TextField
                        style={{ marginTop: 10 }}
                        label="Количество пользователей"
                        type="number"
                        onChange={e => editUserField('count', parseInt(e.target.value))}
                        InputProps={{
                            inputProps: {
                                min: 1
                            }
                        }}
                        fullWidth
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={cancelNewUser} color="secondary">
                        Отмена
                    </Button>
                    <Button onClick={submitNewUser} color="primary">
                        Создать
                    </Button>
                </DialogActions>
            </Dialog>
        )
    }
}

export default observer(CreateUserDialog)
