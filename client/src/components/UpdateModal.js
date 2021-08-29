import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axios from "axios";
import {useAppContext} from "../libs/contextLib.js";
import jwt_decode from "jwt-decode";

export default function UpdateModal({
	item,
}) {
        const { tokenHolder } = useAppContext();
        const [show, setShow] = useState(true);
        const handleClose = () => setShow(false);
	const [actualAmount, setActualAmount] = useState(item.amount);
	const url = "http://spend-it-api.herokuapp.com/v1/graphql"
        const graphqlData= {
                "query":"mutation InsertTransaction($amount: money, $estimate: uuid, $date: timestamptz) { insert_transaction_one(object: {amount: $amount, estimate_id: $estimate, date: $date}){ id } }",
                "variables":{
                },
                "operationName":"InsertTransaction"
        }

	function handleSubmit() {
		console.log("submit", item, actualAmount);
                let headers = {
                        headers: {
                                'Authorization': 'Bearer ' + tokenHolder.getToken(),
                        },
                }

                graphqlData["variables"] = { "amount": actualAmount, "estimate": item.id };

                axios.post(url, graphqlData, headers).then(response=>{
                        console.log(response);
                        if (Object.values(response.data).includes("errors")) {
                                throw("backend error :(");
                        }
                });

                setTimeout(() => {
                        handleClose();
                }, 500);
        }
	
	return (
		<Modal show={show} onHide={handleClose}>
			<Modal.Header>
			<Modal.Title>Update {item.type} Estimate</Modal.Title>
			</Modal.Header>
			<Modal.Body>
				<Form>
 	 				<Form.Group>
    						<Form.Label>Description: <b>{item.description}</b></Form.Label>
						<br/>
    						<Form.Label>Esitmate amount: <b>{item.amount}</b></Form.Label>
						<br/>
    						<Form.Label>Payment on <b>{item.date}</b></Form.Label>
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

