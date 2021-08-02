import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";

export default function UpdateModal({
	item,
}) {

        const [show, setShow] = useState(true);
        const handleClose = () => setShow(false);
	const [actualAmount, setActualAmount] = useState(item.amount);

	function handleSubmit() {
		console.log("submit", item, actualAmount);
	}
	
	return (
		<Modal show={show} onHide={handleClose}>
			<Modal.Header>
			<Modal.Title>Update {item.type} Estimate</Modal.Title>
			</Modal.Header>
			<Modal.Body>
				<Form>
 	 				<Form.Group>
    						<Form.Label>Description: <b>{item.name}</b></Form.Label>
						<br/>
    						<Form.Label>Esitmate amount: <b>{item.amount}</b></Form.Label>
						<br/>
    						<Form.Label>Payment on <b>{item.date.toDateString()}</b></Form.Label>
						<br/>
    						<Form.Label>Actual Amount</Form.Label>
    						<Form.Control type="text" placeholder={item.amount} value={item.actualAmount} onChange={(e)=> setActualAmount(e.target.value)}/>
						
  					</Form.Group>
				</Form>
			</Modal.Body>
			<Modal.Footer>
				<Button variant="secondary" onClick={handleClose}>
					Ask me later
				</Button>
				<Button variant="primary" onClick={handleSubmit}>
					Save
				</Button>
			</Modal.Footer>
		</Modal>
	);
}

