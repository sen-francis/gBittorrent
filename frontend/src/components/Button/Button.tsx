interface ButtonProps {
	buttonText: string,
	onClick: () => void,
	children: React.ReactNode
}

export const Button = (props: ButtonProps) => {
	const {buttonText, onClick, children} = props;
	return 	<div onClick={onClick}>
		{children}
		<div>{buttonText}</div>
	</div>
}
