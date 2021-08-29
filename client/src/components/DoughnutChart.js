import { Doughnut } from 'react-chartjs-2';
import React, { useState, useEffect } from "react";


export default function DoughnutChart({spending, income}) {

	let apiData = {spending: spending, income: income};

	let data = {
  		labels: ["Remaining", "Spending"],
  		datasets: [
    			{
      				data: [Math.max(0, apiData.income - apiData.spending), apiData.spending],
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
		
	};

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
