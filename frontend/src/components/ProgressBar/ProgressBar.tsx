import './ProgressBar.scss';

interface ProgressBarProps {
	progress: number,
	total: number
}

export const ProgressBar = (props: ProgressBarProps) => {
	const {progress, total} = props;

	return <div className="progress-bar">
		<div className="progress-bar__progress" style={{width: `${(progress / total) * 100}%`}}/>
	</div>
}
