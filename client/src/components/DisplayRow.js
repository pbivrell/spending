import React, { useState, useEffect } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import { onError } from "../libs/errorLib";
import "./DisplayRow.css";
import { BsPlus } from "react-icons/bs";
import Button from "react-bootstrap/Button";
import Transaction from "./Transaction";
import TransactionModal from "./TransactionModal";
import {useAppContext} from "../libs/contextLib.js";
import axios from "axios";

import "react-datepicker/dist/react-datepicker.css";


export default function DisplayRow({
	type,
	initalData, 
}) {
	const { tokenHolder } = useAppContext();
	const [loggedIn, setLoggedIn] = useState(tokenHolder.getToken());

	if (!loggedIn) {
		tokenHolder.waitForTokenRefresh().then(() => {
			setLoggedIn(tokenHolder.getToken());
		});
	}

	console.log("Inital data", initalData);
	
	const [show, setShow] = useState(null);
        const handleClose = (t) => {
		console.log("TEST", t, type);
		if (typeof t === "string" && typeof type === "string" && type.toLowerCase() === t.toLowerCase()) {
		}
		setShow(null);
	}
	
        const handleShow = (id, type, description="", amount=0, estimate=true, occurrence=1, period="Day", date=new Date()) => setShow({
		"id": id,
                "type": type,
                "description": description,
                "amount": amount,
                "estimate": estimate,
                "occurrence": occurrence,
                "period": period,
                "date": date,
        });


	const handleClick = () => handleShow("",type);

	function renderNotes() {
		return (
			<div id={type} className="notes">
				<p className="pb-3 mt-4 mb-3 border-bottom">{type}</p>
				<Button onClick={handleClick}>
					<BsPlus size={27} />
				</Button>
				{ 
					initalData === undefined ? <p> Loading</p> : <ListGroup horizontal="md">{renderNotesList(initalData)}</ListGroup>
				}
			</div>
		);
	}

	return (
		<div className="Home">
			{ renderNotes() }
		</div>
	);


	function renderNotesList(notes) {
		console.log("apple", notes);
		return (
			<>
				{notes.map(({ amount, description, id }) => (
					<Transaction onclick={() => handleShow(id, type, description, amount, true)} description={description} amount={amount}/>
				))}
				{ show ? <TransactionModal handleClose={handleClose} data={show}/> : null }
			</>
		);
	}

}

