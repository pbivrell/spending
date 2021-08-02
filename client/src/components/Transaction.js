import ListGroup from "react-bootstrap/ListGroup";


export default function Transaction({
	onclick,
	description,
	amount,
}) {
	return (
		<ListGroup.Item action onClick={onclick}>
			<div>
				<div>{description}</div>
				<div>{amount}</div>
			</div>
		</ListGroup.Item>
	);

}

