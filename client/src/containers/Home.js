import React, { useState, useEffect } from "react";
import "./Home.css";
import DisplayRow from "../components/DisplayRow";
import UpdateModal from "../components/UpdateModal";
import DoughnutChart from "../components/DoughnutChart";
import CategoryBarChart from "../components/CategoryBarChart";
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

export default function Home() {
	
	const [update, setUpdate] = useState([]);

        useEffect(() => {
                loadUpdated();
        }, []);

	async function loadUpdated() {
                let newData = [
                        {
                                amount: 500,
                                name: "test",
                                type: "Spending",
                                date: new Date(),
                        },
                        {
                                amount: 502,
                                name: "test",
                                type: "Income",
                                date: new Date(),
                        },
                        {
                                amount: 503,
                                name: "test",
                                type: "Goal",
                                date: new Date(),
                        },
                ]

                console.log("data", newData);
                setUpdate(newData);
        }

	return (
		<div className="Home">
			<Container id="charts">
  				<Row>
    					<Col>
						<DoughnutChart/>
					</Col>
    					<Col>	
						<h3> Income </h3>
						<CategoryBarChart income={true}/>
						<h3> Spending </h3>
						<CategoryBarChart income={false}/>
					</Col>
  				</Row>
			</Container>
			<DisplayRow type="Spending"/>
			<DisplayRow type="Income"/>
			<DisplayRow type="Goals"/>
				{
					update.map(item => <UpdateModal item={item}/>)
				}
		</div>
	);
}

