import {useState} from 'react';
import './App.css';

import {SelectFile} from '../wailsjs/go/services/FileUploadService'
import {services} from '../wailsjs/go/models'
import {ScrapeTracker} from '../wailsjs/go/services/TrackerService'
const UploadButton = () => {
	const [torrentData, setTorrentData] = useState<services.FileUploadResponse>();
	const [scrapeTrackerData, setScrapeTrackerData] = useState<services.TrackerScrapeResponse>();
	const selectTorrentFile = () => {
		SelectFile().then((data) => setTorrentData(data));		
	}

	const scrapeTracker = () => {
		if (torrentData?.TorrentMetainfo) {	
			ScrapeTracker(torrentData.TorrentMetainfo).then((data) => setScrapeTrackerData(data));
		}	

	}

	return <div>
		<button onClick={() => selectTorrentFile()}>Parse torrent file</button>
		{torrentData && <div> 
			{torrentData.TorrentMetainfo.Announce}
			{torrentData.Err}
		</div>}

		<button onClick={() => scrapeTracker()}>Scrape Tracker</button>

		{scrapeTrackerData && <div>
			{scrapeTrackerData.Seeders}	

		</div>}
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
