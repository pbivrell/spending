import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import DatePicker from "react-datepicker";

export default function TransactionModal({
	handleClose,
	data,
}) {

	const [inData, setData] = useState(data);

	function onChange(e, type) {
		

		let value =  e;

		if (type !== "date") {
			value = e.target.value;
		}


		if (type === "amount" || type === "occurance" ) {
			value = parseInt(value);
			if (isNaN(value)) {
				value = 0;
			}
		}

		let newData = {
			id: inData.id,
			amount: inData.amount,
			name: inData.name,
			type: inData.type,
			update: inData.update,
			occurance: inData.occurance,
			period: inData.period,
			date: inData.date,
		}
		console.log(newData);
		newData[type] = value;
		setData(newData);
	}

	function handleSubmit() {
		console.log(inData);
		handleClose();
	}

	function nextOccurances(indate, occurance, period) {
		var date = new Date(indate);
		console.log(date, occurance, period);
		if (period === "Day") {
			date.setDate(date.getDate() + occurance);
		} else if (period === "Week") {
			date.setDate(date.getDate() + (7 * occurance));
		} else if (period === "Month") {
			date.setMonth(date.getMonth() + occurance);
		} else if (period === "Year") {
			date.setFullYear(date.getFullYear() + occurance);
		}
		console.log(date);
		return date;
	}
	

	return (
		<Modal show={true} onHide={handleClose}>
			<Modal.Header closeButton>
				<Modal.Title>New {data.type}</Modal.Title>
			</Modal.Header>
			<Modal.Body>
				<Form>
 	 				<Form.Group>
    						<Form.Label>Amount</Form.Label>
    						<Form.Control type="text" placeholder="" value={inData.amount} onChange={(e)=>onChange(e, "amount")}/>
    						<Form.Label>Description</Form.Label>
    						<Form.Control type="text" placeholder="" value={inData.name} onChange={(e)=>onChange(e, "name")}/>
						<Form.Label>Reoccurance</Form.Label>
						<br/><span>Every</span> <Form.Control type="text" placeholder="" value={inData.occurance} onChange={(e)=>onChange(e, "occurance")}/>
    						<Form.Control as="select" onChange={(e)=>onChange(e, "period")}>
      							<option>Day</option>
      							<option>Week</option>
      							<option>Month</option>
      							<option>Year</option>
    						</Form.Control>
						<span>Start </span><DatePicker selected={inData.date} onChange={(date)=>onChange(date, "date")} />
						<br/>
						<Form.Label>Estimate</Form.Label>
						<Form.Check type="checkbox" label={`Checking this box means the amount you entered is an esitimate. Every ${inData.occurance} ${inData.period} you want to save the real value`}/>
						<em> The next time this date will happen is { nextOccurances(inData.date, inData.occurance, inData.period).toDateString() }</em>
  					</Form.Group>
				</Form>
			</Modal.Body>
			<Modal.Footer>
				<Button variant="secondary" onClick={handleClose}>
					Close
				</Button>
				{ inData.update ? <Button variant="danger" onClick={handleClose}>Delete</Button>: null }
				<Button variant="primary" onClick={handleSubmit}>
					Save Changes
				</Button>
			</Modal.Footer>
		</Modal>
	);

}

