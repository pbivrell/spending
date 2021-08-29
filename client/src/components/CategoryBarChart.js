import React from 'react';
import { Bar } from 'react-chartjs-2';
import { useState, useEffect } from "react";
import 'chartjs-plugin-labels';


const options = {
  scales: {
    yAxes: [
      {
        ticks: {
          beginAtZero: true,
        },
      },
    ],
  },
};

export default function CategoryBarChart ( {

	income,
	inData,
}) {

	const [data, setData] = useState([]);

	async function getData() {
	
		let apiData = {}

		let colors = [
      			'rgba(40, 179, 23, 0.4)',
      			'rgba(40, 179, 23, 1)',
		]

		if (!income) {
			colors = [
				'rgba(252, 3, 61, 0.2)',
				'rgba(252, 3, 61, 1.0)',
			]
		}

		let graphData = {
			labels: [],
			datasets: [
				{
					label: 'Estimate',
					data: [],
      					backgroundColor: colors[0],
				},
				{
      					label: 'Actual',
      					data: [],
      					backgroundColor: colors[1],
				},
			],
		}

		if (inData) {
		Object.values(inData).forEach((item) => {
			graphData.labels.push(item.description);
			graphData.datasets[0].data.push(parseFloat(item.amount.substring(1)));
			let actual = 0;
			if(apiData[item.description]) {
				actual = parseFloat(item.amount.substring(1));
			}
			graphData.datasets[1].data.push(actual);
		});
		}

		setData(graphData);
	}

	useEffect(() => {
        	getData();
    	}, [inData]);

	return (
		<Bar data={data} options={options} />
   	);
}

