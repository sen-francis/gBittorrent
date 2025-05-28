import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';

import {SelectFile} from '../wailsjs/go/services/FileUploadService'

const UploadButton = () => {
	const selectTorrentFile = () => {
		SelectFile().then((data) => console.log(data));		
	}


	return <div>
		<button onClick={() => selectTorrentFile()}>Parse torrent file</button>
	</div>
}

function App() {
	return (
		<div id="App">
			<UploadButton/>
		</div>
	)
}

export default App
