import './Client.css';
import {useEffect, useState, useRef} from 'react';
import {SelectFile} from '../../wailsjs/go/services/FileUploadService'
import {DownloadTorrent} from '../../wailsjs/go/services/TorrentService'
import {services, torrent} from '../../wailsjs/go/models'
import {ScrapeTracker} from '../../wailsjs/go/services/TrackerService'
import { Modal } from '../components/Modal';

const BYTES_IN_GB = 1000000000;
const BYTES_IN_MB = 1000000;
const BYTES_IN_KB = 1000;

const roundToHundreths = (num: number) => {
	return Math.round(100 * num) / 100;
}

const getFormattedFileSizeString = (bytes: number) => {
	if (bytes >= BYTES_IN_GB) {
		return roundToHundreths(bytes / BYTES_IN_GB) + " GB";
	}
	else if (bytes >= BYTES_IN_MB) {
		return roundToHundreths(bytes / BYTES_IN_MB) + " MB";
	}
	else if (bytes >= BYTES_IN_KB) {
		return roundToHundreths(bytes / BYTES_IN_KB) + " KB";
	}	
	return bytes + " B";
}

const TorrentInformation = (props: {torrentMetainfo: torrent.TorrentMetainfo, scrapeTrackerData: services.TrackerScrapeResponse}) => {
	const {torrentMetainfo, scrapeTrackerData} = props;

	return <div className='torrent-information'>
		<span><b>Torrent information</b></span>
		<div className='torrent-information__body'>
			<span className='torrent-information__field'><b>Name:</b> {scrapeTrackerData.Name}</span>
			<span className='torrent-information__field'><b>Size:</b> {getFormattedFileSizeString(torrentMetainfo.Size)}</span>
			<span className='torrent-information__field'><b>Downloaded:</b> {scrapeTrackerData.Downloaded}</span>
			<span className='torrent-information__field'><b>Seeders:</b> {scrapeTrackerData.Seeders}</span>
			<span className='torrent-information__field'><b>Leechers:</b> {scrapeTrackerData.Leechers}</span>
			<span className='torrent-information__field'><b>Creation date</b>: {torrentMetainfo.CreationDate}</span>
			<span className='torrent-information__field'><b>Created by</b>: {torrentMetainfo.CreatedBy}</span>
			{torrentMetainfo.Comment && <span className='torrent-information__field'><b>Comment:</b> {torrentMetainfo.Comment}</span>}
		</div>
	</div>
}

const TorrentFileInformation = (props: {torrentInfo: torrent.TorrentInfo}) => {
	const {torrentInfo} = props;
	return <table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Size</th>
			</tr>
		</thead>
		<tbody>
			{torrentInfo.FileInfoList.map((fileInfo: torrent.FileInfo, index) => {
				return <tr key={index}>
					<td>{fileInfo.Path}</td>
					<td>{getFormattedFileSizeString(fileInfo.Length)}</td>
				</tr>	
			})}
		</tbody>
	</table>
}

const Menu = () => {
	const [torrentData, setTorrentData] = useState<services.FileUploadResponse>();
	const [scrapeTrackerData, setScrapeTrackerData] = useState<services.TrackerScrapeResponse>();
	const [isModalOpen, setIsModalOpen] = useState<boolean>(false);	

	const addTorrent = () => {
		SelectFile().then((data) => {
			setTorrentData(data);
			scrapeTracker(data);
		});
	}

	const scrapeTracker = (data: services.FileUploadResponse) => {
		if (data?.TorrentMetainfo) {	
			ScrapeTracker(data.TorrentMetainfo).then((data) => {
				setScrapeTrackerData(data);
				setIsModalOpen(true);
			});
		}
	}

	const onModalClose = () => {
		setScrapeTrackerData({} as services.TrackerScrapeResponse);
		setTorrentData({} as services.FileUploadResponse);
		setIsModalOpen(false);
	}

	const onModalSubmit = () => {
		if (!torrentData?.TorrentMetainfo) {
			return;	
		}
		DownloadTorrent(torrentData.TorrentMetainfo);
	}

	return <div>
		<button onClick={() => addTorrent()}>Add</button>
		<button>Remove</button>
		<Modal onClose={onModalClose} isModalOpen={isModalOpen} submitText="Download" onSubmit={onModalSubmit}>
			{(torrentData?.TorrentMetainfo && scrapeTrackerData != null) && <>
				<TorrentInformation torrentMetainfo={torrentData.TorrentMetainfo} scrapeTrackerData={scrapeTrackerData} />
				<TorrentFileInformation torrentInfo={torrentData.TorrentMetainfo.Info}/>
			</>}
		</Modal>
	</div>
}

const Torrent = () => {
	return <tr>

	</tr>
}

const TorrentList = () => {
	return <table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Size</th>
				<th>Status</th>
				<th>ETA</th>
				<th>Speed (Down/Up)</th>
				<th>Peers</th>
				<th>Seeds</th>
			</tr>	
		</thead>
		<tbody>
		</tbody>
	</table>
}

export const Client = () => {
	return <div>
		<Menu/>	
		<TorrentList/>
	</div>
}
