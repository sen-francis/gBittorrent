import './ProgressBar.scss';

interface ProgressBarProps {
	progress: number,
	total: number
}

export const ProgressBar = (props: ProgressBarProps) => {
	const {progress, total} = props;
	const isComplete = progress == total;
	let progressBarClasses = ["progress-bar"]
	if (isComplete) {
		progressBarClasses.push("progress-bar--complete");
	}
	return <div className={progressBarClasses.join(" ")}>
		{!isComplete && <div className="progress-bar__progress" style={{width: `${(progress / total) * 100}%`}}/>}
	</div>
}
