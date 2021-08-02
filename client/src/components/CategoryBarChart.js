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

export default function CategoryBarChart (
	income=true,
) {

	const [data, setData] = useState([]);

	async function getData() {
	
		let apiData = [
			{
				name: "other", 
				estimate: 50,
				actual: 200,
			},{
				name: "name", 
				estimate: 92,
				actual: 81,
			},{
				name: "bear", 
				estimate: 1,
				actual: 21,
			},
		];

		let colors = [
      			'rgba(40, 179, 23, 0.4)',
      			'rgba(40, 179, 23, 1)',
		]

		console.log(income.income)
			
		if (!income.income) {
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

		apiData.forEach((item) => {
			graphData.labels.push(item.name);
			graphData.datasets[0].data.push(item.estimate);
			graphData.datasets[1].data.push(item.actual);
		});
	
		console.log(graphData);
		setData(graphData);
	}

	useEffect(() => {
        	getData();
    	}, []);

	return (
		<Bar data={data} options={options} />
   	);
}

