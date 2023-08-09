import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import License from 'containers/License'
import { withSnackbar } from 'notistack';
import BugReports from 'containers/BugReports';

class App extends React.Component {
	state = {
		page: 0
	}

	constructor(props) {
        super(props)
        window.notify = (message, variant, options) => {
            if (!options) {
                options = {}
            }
            options.variant = variant ? variant : 'info'
            props.enqueueSnackbar(message, options);
        }
		window.closeNotify = props.closeSnackbar
    }

	render() {
		const { page } = this.state
		return (
			<>
				<AppBar color="default" position="relative">
					<Tabs value={page} onChange={(_, page) => this.setState({ page })}>
						<Tab label="Пользователи" />
						<Tab label="Баг-репорты" />
					</Tabs>
				</AppBar>
				{(() => {
					switch (page) {
						case 0: return <License />
						case 1: return <BugReports />
					}
				})()}

			</>
		)
	}
}

export default withSnackbar(App)
