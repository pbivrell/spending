import React, { useState, useEffect } from "react";
import "./Home.css";
import DisplayRow from "../components/DisplayRow";
import UpdateModal from "../components/UpdateModal";
import DoughnutChart from "../components/DoughnutChart";
import CategoryBarChart from "../components/CategoryBarChart";
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import axios from "axios";
import {useAppContext} from "../libs/contextLib.js";
import { add } from 'date-fns'


export default function Home() {
	
	const {tokenHolder} = useAppContext();

	const [data, setData] = useState({});

	const [token, setToken] = useState(tokenHolder.getToken());

	if (token === null)  {
		tokenHolder.waitForTokenRefresh().then(() => setToken(tokenHolder.getToken()));
	}

	const query = `
query Updates{
  estimate(where: {deleted: {_is_null: true}}){
    id
    description
    amount
    created
    date
    deleted
    estimate
    occurrence
    owner_id
    period
    type
    transactions(order_by: {date: asc_nulls_last}, limit: 1) {
      date
    }
  }
}`
	const url = "http://spend-it-api.herokuapp.com/v1/graphql"
        const graphqlData= {
                "query":query,
                "variables":{},
                "operationName":"Updates"
        }

	function getGraphql() {
		if (token === null) {
			return
		}
                let headers = {
                        headers: {
                                'Authorization': 'Bearer ' +token,
                        },
                }

                axios.post(url, graphqlData, headers).then(response=>{
                        if (Object.values(response.data).includes("errors")) {
                                throw("backend error :(");
                        }
			setData(sortData(response.data.data.estimate));
                });
	}

	function sortData(data) {
		let categorized  = {
			totalIncome: 0,
			totalSpending: 0,
			update: [],
		};

		data.forEach((item) => {
			if (item.type === "spending"){
				categorized["totalSpending"] += parseFloat(item.amount.substring(1));
			}else if(item.type === "income" ) {
				categorized["totalIncome"] += parseFloat(item.amount.substring(1));
			}
			if (!categorized[item.type]) {
				categorized[item.type] = [item];
			}else {
				categorized[item.type].push(item);
			}

			console.log("THESE ARE", item, item.date) 

			if(item.date) {
				let searchDate = item.date;
				if(item.transactions.length > 0) {
					searchDate = item.transactions[0].date;
				}
				for (let n of lastOccurance(searchDate, item.occurrence, item.period)) {
					console.log("pear", n);
					categorized.update.unshift({
						description: item.description,
						amount: item.amount,
						id: item.id,
						date: n.toString(),
					});
				}
			}
		});


		return categorized;
	}


	function* lastOccurance(date, occurrence, period) {
                
  		let now = new Date();
 
		let d = new Date(date);

  		let steps = 1;
      
		console.log("Generating", now, d, occurrence, period)

  		while(true) {
  
    			let amount = { years: steps * occurrence};
  
    			if (period === "day") {
 		  		amount = {days: steps * occurrence}	
    			} else if (period === "week") {
  	  			amount = {weeks: steps *  occurrence}
    			} else if (period === "month") {
  	  			amount = {months: steps * occurrence}
    			} 
  	
  			let iterDate = add(d, amount);

			if (steps > 10){
				break;
			}

    			if(iterDate > now) {
    				break;
    			}

			console.log("JELLO", iterDate, amount, steps)
    			yield iterDate;
			steps++
		}

		console.log("WE GENERaTED");
	}


        useEffect(() => {

		if ( Object.keys(data).length === 0) {
                	getGraphql();
		}
        }, [token]);

	console.log("Spending", data);

	return (
		<div className="Home">
			<Container id="charts">
  				<Row>
    					<Col>
						<DoughnutChart spending={data.totalSpending} income={data.totalIncome}/>
					</Col>
    					<Col>	
						<h3> Income </h3>
						<CategoryBarChart income={true} inData={data.income}/>
						<h3> Spending </h3>
						<CategoryBarChart income={false} inData={data.spending}/>
					</Col>
  				</Row>
			</Container>
			<DisplayRow type="Spending" initalData={data.spending}/>
			<DisplayRow type="Income" initalData={data.income}/>
			<DisplayRow type="Goals" initalData={data.goals}/>
				{
					data.update ? data.update.map(item => <UpdateModal item={item}/>) : <></>
					
				}
		</div>
	);
}

