const efficiencyChart = document.getElementById('efficiencyChart');
const evaporationChart = document.getElementById('evaporationChart');
const data = window.Stats;
console.log(JSON.stringify(data));

const efficiencyData = data.filter(row => row.Efficiency !== null && !isNaN(row.Efficiency));
new Chart(efficiencyChart, {
    type: 'bar',
    data: {
        labels: efficiencyData.map(row => row.FinishedTimeString),
        datasets: [
            {
                label: 'Efficiency',
                data: efficiencyData.map(row => row.Efficiency),
                customText: efficiencyData.map(row => row.RecipeName)
            }
        ]
    },
    options: {
        plugins: {
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const index = context.dataIndex;
                        const dataset = context.dataset;
                        const value = dataset.data[index];
                        const note = dataset.customText ? dataset.customText[index] : '';
                        return `Value: ${value} — ${note}`;
                    }
                }
            }
        }
    }
});
const efficiencyFiltered = data.filter(v => v.Efficiency !== null && !isNaN(v.Efficiency));
const efficiencyAvg = efficiencyFiltered.reduce((sum, d) => sum + d.Efficiency, 0) / efficiencyFiltered.length;
document.getElementById('efficiencyAvg').textContent = efficiencyAvg.toFixed(1);

const evaporationFiltered = data.filter(v => v.Evaporation !== null && !isNaN(v.Evaporation));
const evaporationAvg = evaporationFiltered.reduce((sum, d) => sum + d.Evaporation, 0) / evaporationFiltered.length;
document.getElementById('evaporationAvg').textContent = evaporationAvg.toFixed(1);


const evaporationData = data.filter(row => row.Evaporation !== null && !isNaN(row.Evaporation));
new Chart(evaporationChart, {
    type: 'bar',
    data: {
        labels: evaporationData.map(row => row.FinishedTimeString),
        datasets: [
            {
                label: 'Evaporation',
                data: evaporationData.map(row => row.Evaporation),
                customText: evaporationData.map(row => row.RecipeName)
            }
        ]
    },
    options: {
        plugins: {
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const index = context.dataIndex;
                        const dataset = context.dataset;
                        const value = dataset.data[index];
                        const note = dataset.customText ? dataset.customText[index] : '';
                        return `Value: ${value} — ${note}`;
                    }
                }
            }
        }
    }
});