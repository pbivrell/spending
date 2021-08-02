import { Doughnut } from 'react-chartjs-2';
import React, { useState, useEffect } from "react";


export default function DoughnutChart() {

	const [data, setData] = useState([]);

	async function getData() {
		
		let apiData = {spending: 1504, income: 2000};

		setData({
  			labels: ["Remaining", "Spending"],
  			datasets: [
    			{
      				data: [apiData.income - apiData.spending, apiData.spending],
      				backgroundColor: [
        				'rgba(66, 135, 245, 0.7)',
        				'rgba(252, 3, 61, 0.7)',
      				],
      				borderColor: [
        				'rgba(66, 135, 245, 0.1)',
        				'rgba(252, 3, 61, 0.1)',
      				],
      				borderWidth: 2,
    			},
  			],
		})
	};

	useEffect(() => {
    	getData();
    }, []);

	const options = {
    		plugins: {
      			labels: {
        			render: 'value',
			}
		}
	}

	return (
    		<Doughnut data={data} options={options}/>
	);
}
