import {useState} from 'react';
import './App.css';

import {SelectFile} from '../wailsjs/go/services/FileUploadService'
import {services} from '../wailsjs/go/models'
import {ScrapeTracker} from '../wailsjs/go/'
const UploadButton = () => {
	const [torrentData, setTorrentData] = useState<services.FileUploadResponse>();
	const selectTorrentFile = () => {
		SelectFile().then((data) => setTorrentData(data));		
	}
	
	return <div>
		<button onClick={() => selectTorrentFile()}>Parse torrent file</button>
		{torrentData && <div> 
			{torrentData.TorrentMetainfo.Announce}
			{torrentData.Err}
		</div>}
	
		<Dialog>
			</Dialog>

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
