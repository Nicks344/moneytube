import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Button from '@material-ui/core/Button'
import { SnackbarProvider } from 'notistack';

const theme = createMuiTheme({
	palette: {
		type: 'dark',
		primary: {
			light: '#f85300',
			dark: '#ff7f08',
			main: '#ff6e0a'
		},
		error: {
			main: '#ff2929'
		},
		success: {
			light: '#309308',
			dark: '#8fd513',
			main: '#8fd513'
		},
		secondary: {
			light: '#000',
			dark: '#fff',
			main: '#fff'
		},
	}
})

ReactDOM.render(
	<ThemeProvider theme={theme}>
		<CssBaseline />
		<SnackbarProvider
			maxSnack={1}
			anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
			action={key => (
				<Button color="inherit" onClick={() => window.closeNotify(key)}>
					{'Close'}
				</Button>
			)}>
			<App />
		</SnackbarProvider>
	</ThemeProvider>,
	document.getElementById('root')
);
