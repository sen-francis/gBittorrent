import { ProgressBar } from '../components/ProgressBar/ProgressBar';
import './Client.scss';
import { Toolbar } from './Toolbar/Toolbar';
import Folder from "../assets/images/folder.svg?react"
import File from "../assets/images/file.svg?react"
import { createContext, useContext, useState } from 'react';
import { services } from '../../wailsjs/go/models';
import Trash from "../assets/images/trash.svg?react"
import {EventsEmit} from '../../wailsjs/runtime/runtime';
import { Button } from '../components/Button/Button';

const Torrent = (props: {torrentData: services.FileUploadResponse}) => {
	const {torrentData} = props;
	const total = 100;
	const progress = 50;
	const progressText = `${progress} of ${total} (${progress / total * 100}%)`;
	const torrentInfo = torrentData.TorrentMetainfo.Info;
	const singleFileMode = torrentInfo.FileInfoList.length == 1;
	const getTorrentName = () => {
		const torrentInfo = torrentData.TorrentMetainfo.Info;
		if (singleFileMode) {
			return torrentInfo.FileInfoList[0].Path.join("");	
		}
		return torrentInfo.DirectoryName;
	}


	const stopTorrentDownload = () => {
		if (!torrentData?.TorrentMetainfo?.InfoHashStr) {
			return;	
		}
		EventsEmit(torrentData.TorrentMetainfo.InfoHashStr)
	}


	return <div className='torrent'>
		{singleFileMode ? <File/> : <Folder/>}
		<div className="torrent_info">
			<div>{getTorrentName()}</div>
			<div className='torrent_info_details'>{progressText}</div>
			<div className='torrent_info_progress-container'>
				<ProgressBar total={100} progress={50} />
				<Button onClick={() => stopTorrentDownload()}>
					<Trash/>
				</Button>
			</div>
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
