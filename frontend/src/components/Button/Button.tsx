import './Button.scss';

interface ButtonProps {
	buttonText: string,
	onClick: () => void,
	children: React.ReactNode
}

export const Button = (props: ButtonProps) => {
	const {buttonText, onClick, children} = props;
	return 	<div className="button" onClick={onClick}>
		{children}
		<div>{buttonText}</div>
	</div>
}
