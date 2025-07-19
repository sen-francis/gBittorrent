import { ProgressBar } from '../components/ProgressBar/ProgressBar';
import './Client.scss';
import { Toolbar } from './Toolbar/Toolbar';
import Folder from "../assets/images/folder.svg?react"
import { createContext, useContext, useState } from 'react';
import { services } from '../../wailsjs/go/models';

const Torrent = (props: {torrentData: services.FileUploadResponse}) => {
	const {torrentData} = props;
	const total = 100;
	const progress = 50;
	const progressText = `${progress} of ${total} (${progress / total * 100}%)`;

	const getTorrentName = () => {
		const torrentInfo = torrentData.TorrentMetainfo.Info;
		if (torrentInfo.FileInfoList.length == 1) {
			return torrentInfo.FileInfoList[0].Path.join("");	
		}
		return torrentInfo.DirectoryName;
	}

	return <div className='torrent'>
		<Folder/>
		<div className="torrent_info">
			<div>{getTorrentName()}</div>
			<div className='torrent_info_details'>{progressText}</div>
			<ProgressBar total={100} progress={50} />
			<div className='torrent_info_details'>speed: </div>
		</div>
	</div>
}

const TorrentList = () => {
	const {torrentList} = useContext(ClientContext);
	return <div className='torrent-container'>
		{torrentList.map((torrent) => {
			return <Torrent torrentData={torrent}/>;
		})}
	</div>
}

interface ClientContextType {
	torrentList: Array<services.FileUploadResponse>;
	setTorrentList: React.Dispatch<React.SetStateAction<Array<services.FileUploadResponse>>>;
}

export const ClientContext = createContext<ClientContextType>({} as ClientContextType);

export const Client = () => {
	const [torrentList, setTorrentList] = useState<Array<services.FileUploadResponse>>([]);
	return <div className='client'>
		<ClientContext.Provider value={{torrentList: torrentList, setTorrentList: setTorrentList}}>
			<Toolbar />
			<TorrentList/>
		</ClientContext.Provider>
	</div>
}
