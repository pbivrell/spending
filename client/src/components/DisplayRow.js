import React, { useState, useEffect } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import { onError } from "../libs/errorLib";
import "./DisplayRow.css";
import { BsPlus } from "react-icons/bs";
import Button from "react-bootstrap/Button";
import Transaction from "./Transaction";
import TransactionModal from "./TransactionModal";

import "react-datepicker/dist/react-datepicker.css";


export default function DisplayRow({
	type,
}) {
	const [notes, setNotes] = useState([]);
	//const { isAuthenticated } = useAppContext();
	const [isLoading, setIsLoading] = useState(true);
	
	const [show, setShow] = useState(null);
        const handleClose = () => setShow(null);
        const handleShow = (id, type, name="", amount=0, update=false, occurance=1, period="Day", date=new Date()) => setShow({
		"id": id,
                "type": type,
                "name": name,
                "amount": amount,
                "update": update,
                "occurance": occurance,
                "period": period,
                "date": date,
        });


	const handleClick = () => handleShow(type);

	useEffect(() => {
		async function onLoad() {

			try {
				const notes = await loadNotes();
				setNotes(notes);
			} catch (e) {
				onError(e);
			}

			setIsLoading(false);
		}

		onLoad();
	}, []);
	
	function loadNotes() {
		return [{Id: 0, Amount: 500, Description: "Name"},{Id: 1, Amount: 500, Description: "Other"} ];
	}


	function renderNotes() {
		return (
			<div id={type} className="notes">
				<p className="pb-3 mt-4 mb-3 border-bottom">{type}</p>
				<Button onClick={handleClick}>
					<BsPlus size={27} />
				</Button>
				<ListGroup horizontal="md">{!isLoading && renderNotesList(notes)}</ListGroup>
			</div>
		);
	}

	return (
		<div className="Home">
			{ renderNotes() }
		</div>
	);


	function renderNotesList(notes) {
		return (
			<>
				{notes.map(({ Amount, Description, Id }) => (
					<Transaction onclick={() => handleShow(Id, type, Description, Amount, true)} description={Description} amount={Amount}/>
				))}
				{ show ? <TransactionModal handleClose={handleClose} data={show}/> : null }
			</>
		);
	}

}

