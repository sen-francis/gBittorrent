import { ProgressBar } from '../components/ProgressBar/ProgressBar';
import './Client.scss';
import { Toolbar } from './Toolbar/Toolbar';
import Folder from "../assets/images/folder.svg?react"

interface TorrentProps {
	name: string,
	
}
const Torrent = () => {
	const name = "ubuntu";
	const total = 100;
	const progress = 50;
	const progressText = `${progress} of ${total} (${progress / total * 100}%)`;
	return <div className='torrent'>
		<Folder/>
		<div className="torrent_info">
			<div>{name}</div>
			<div className='torrent_info_details'>{progressText}</div>
			<ProgressBar total={100} progress={50} />
			<div className='torrent_info_details'>speed: </div>
		</div>
	</div>
}

const TorrentList = () => {
	return <div className='torrent-container'>
		<Torrent/>
		<Torrent/>
		<Torrent/>
	</div>
}

export const Client = () => {
	return <div className='client'>
		<Toolbar />
		<TorrentList/>
	</div>
}
