import './Modal.css';
import {ReactNode, useEffect, useRef} from 'react';

interface ModalProps {
	onClose: () => void;
	isModalOpen: boolean;
	children: ReactNode
	onSubmit: () => void;
	submitText: string;
}

export const Modal = (props: ModalProps) => {
	const {onClose, isModalOpen, children, onSubmit, submitText} = props;
	const modalRef = useRef<HTMLDialogElement>(null);
	useEffect(() => {
		if (!modalRef.current) {
			return;	
		}
		if (isModalOpen) {
			modalRef.current.showModal();	
		}
		else {
			modalRef.current.close();
			onClose();
		}
	}, [isModalOpen]);

	return <dialog ref={modalRef}>
		<div className="modal__body">
			{children}		
		</div>
		<div className="modal__footer">
			<button onClick={onSubmit}>{submitText}</button>
			<button onClick={onClose}>Cancel</button>
		</div>
	</dialog>
}
