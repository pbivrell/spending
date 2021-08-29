import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import DatePicker from "react-datepicker";
import axios from "axios";
import {useAppContext} from "../libs/contextLib.js";
import jwt_decode from "jwt-decode";

export default function TransactionModal({
	handleClose,
	data,
}) {

	const [inData, setData] = useState(data);
	const { tokenHolder } = useAppContext();


	const url = "http://spend-it-api.herokuapp.com/v1/graphql"
        const graphqlData= {
                "query":"mutation InsertEstimate($object: estimate_insert_input = {}) { insert_estimate_one(object: $object){ id } }",
                "variables":{
                },
                "operationName":"InsertEstimate"
        }


	function onChange(e, type) {
		

		let value =  e;

		if (type !== "date") {
			value = e.target.value;
		}


		if (type === "amount" || type === "occurrence" ) {
			value = parseInt(value);
			if (isNaN(value)) {
				value = 0;
			}
		}

		let owner = "";
		if (tokenHolder !== null && tokenHolder.getToken() !== null ) {
			owner = jwt_decode(tokenHolder.getToken())
		}

		let newData = {
			amount: inData.amount,
			description: inData.description,
			type: inData.type.toLowerCase(),
			owner_id: "owner",
			estimate: inData.estimate,
			occurrence: (inData.estimate) ? inData.occurrence : null,
			period: (inData.estimate) ? inData.period.toLowerCase() : null,
			date: (inData.estimate) ? inData.date : null,
		}
		console.log(newData);
		newData[type] = value;
		setData(newData);
	}

	function handleSubmit() {
		let headers = {
  			headers: {
      				'Authorization': 'Bearer ' + tokenHolder.getToken(),
			},
		}
	
		graphqlData["variables"] = { object: inData};

	        axios.post(url, graphqlData, headers).then(response=>{
			console.log(response);
			if (Object.values(response.data).includes("errors")) {
				throw("backend error :(");
			}
		});

		setTimeout(() => {
			handleClose(inData.type);
		}, 2000);
	}

	function nextOccurrences(indate, occurrence, period) {
		var date = new Date(indate);
		console.log(date, occurrence, period);
		if (period === "Day") {
			date.setDate(date.getDate() + occurrence);
		} else if (period === "Week") {
			date.setDate(date.getDate() + (7 * occurrence));
		} else if (period === "Month") {
			date.setMonth(date.getMonth() + occurrence);
		} else if (period === "Year") {
			date.setFullYear(date.getFullYear() + occurrence);
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
    						<Form.Control type="text" placeholder="" value={inData.description} onChange={(e)=>onChange(e, "description")}/>
						<Form.Label>Reoccurrence</Form.Label>
						<br/><span>Every</span> <Form.Control type="text" placeholder="" value={inData.occurrence} onChange={(e)=>onChange(e, "occurrence")}/>
    						<Form.Control as="select" onChange={(e)=>onChange(e, "period")}>
      							<option>Day</option>
      							<option>Week</option>
      							<option>Month</option>
      							<option>Year</option>
    						</Form.Control>
						<span>Start </span><DatePicker selected={inData.date} onChange={(date)=>onChange(date, "date")} />
						<br/>
						<Form.Label>Estimate</Form.Label>
						<Form.Check type="checkbox" label={`Checking this box means the amount you entered is an esitimate. Every ${inData.occurrence} ${inData.period} you want to save the real value`} onChange={(e)=>onChange(e.target.checked, "estimate")}/>
						<em> The next time this date will happen is { nextOccurrences(inData.date, inData.occurrence, inData.period).toDateString() }</em>
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

