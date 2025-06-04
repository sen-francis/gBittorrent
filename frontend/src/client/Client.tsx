import './Client.css';
import {useState} from 'react';
import {SelectFile} from '../../wailsjs/go/services/FileUploadService'
import {services, torrent} from '../../wailsjs/go/models'
import {ScrapeTracker} from '../../wailsjs/go/services/TrackerService'

const TorrentInformation = (props: {torrentMetainfo: torrent.TorrentMetainfo, scrapeTrackerData: services.TrackerScrapeResponse}) => {
	const {torrentMetainfo, scrapeTrackerData} = props;

	return <div>
		<span>Torrent information</span>
		<div>
			<span>Name: {scrapeTrackerData.Name}</span>
			<span>Announce: {torrentMetainfo.Announce}</span>
			<span>Info hash: {torrentMetainfo.InfoHash}</span>
			<span>Size: {torrentMetainfo.Size}</span>
			<span>Downloaded: {scrapeTrackerData.Downloaded}</span>
			<span>Seeders: {scrapeTrackerData.Seeders}</span>
			<span>Leechers: {scrapeTrackerData.Leechers}</span>
			<span>Creation date: {torrentMetainfo.CreationDate}</span>
			<span>Created by: {torrentMetainfo.CreatedBy}</span>
			{torrentMetainfo.Comment && <span>Comment: {torrentMetainfo.Comment}</span>}
		</div>
	</div>
}

const TorrentFileInformation = (props: {torrentInfo: torrent.TorrentInfo}) => {
	const {torrentInfo} = props;
	return <table>
		<tr>
			<td>Name</td>
			<td>Size</td>
		</tr>
		{torrentInfo.FileInfoList.map((fileInfo: torrent.FileInfo) => {
			return <tr>
				<td>{fileInfo.Path}</td>
				<td>{fileInfo.Length}</td>
			</tr>	
		})}
	</table>
}

interface AddTorentModalProps {
	torrentData: services.FileUploadResponse;
	scrapeTrackerData: services.TrackerScrapeResponse;
	onClose: () => void;
}

const AddTorrentModal = (props: AddTorentModalProps) => {
	const {torrentData, scrapeTrackerData, onClose} = props;

	return <dialog>
		<div>
			<TorrentInformation torrentMetainfo={torrentData.TorrentMetainfo} scrapeTrackerData={scrapeTrackerData} />
			<TorrentFileInformation torrentInfo={torrentData.TorrentMetainfo.Info}/>
		</div>
		<div>
			<button>Download</button>
			<button onClick={onClose}>Cancel</button>
		</div>
	</dialog>
}

const Menu = () => {
	const [torrentData, setTorrentData] = useState<services.FileUploadResponse>();
	const [scrapeTrackerData, setScrapeTrackerData] = useState<services.TrackerScrapeResponse>();
	const [isModalOpen, setIsModalOpen] = useState<boolean>(false);	

	const addTorrent = () => {
		SelectFile().then((data) => {
			setTorrentData(data);
			scrapeTracker();
		});
	}

	const scrapeTracker = () => {
		if (torrentData?.TorrentMetainfo) {	
			ScrapeTracker(torrentData.TorrentMetainfo).then((data) => {
				setScrapeTrackerData(data);
				setIsModalOpen(true);
			});
		}
	}

	const onModalClose = () => {
		setScrapeTrackerData(undefined);
		setTorrentData(undefined);
		setIsModalOpen(false);
	}

	const shouldShowModal = isModalOpen && torrentData != null && scrapeTrackerData != null;

	return <div>
		<button onClick={() => addTorrent()}>Add</button>
		<button>Remove</button>
		{shouldShowModal && <AddTorrentModal torrentData={torrentData} scrapeTrackerData={scrapeTrackerData} onClose={onModalClose} />}
	</div>
}

const Torrent = () => {
	return <tr>

	</tr>
}

const TorrentList = () => {
	return <table>
		<tr>
			<td>Name</td>
			<td>Size</td>
			<td>Status</td>
			<td>ETA</td>
			<td>Speed (Down/Up)</td>
			<td>Peers</td>
			<td>Seeds</td>
		</tr>
	</table>
}

export const Client = () => {
	return <div>
		<Menu/>	
		<TorrentList/>
	</div>
}
