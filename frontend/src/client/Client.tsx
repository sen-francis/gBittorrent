import './Client.scss';
import { Toolbar } from './Toolbar/Toolbar';



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
		<Toolbar />
		<TorrentList/>
	</div>
}
