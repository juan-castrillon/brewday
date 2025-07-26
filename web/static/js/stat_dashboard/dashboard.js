const efficiencyChart = document.getElementById('efficiencyChart');
const evaporationChart = document.getElementById('evaporationChart');
const data = window.Stats;
new Chart(efficiencyChart, {
    type: 'bar',
    data: {
        labels: data.map(row => row.FinishedTimeString),
        datasets: [
            {
                label: 'Efficiency',
                data: data.map(row => row.Efficiency),
                customText: data.map(row => row.RecipeName)
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
const efficiencyAvg = data.reduce((sum, d) => sum + d.Efficiency, 0) / data.length;
document.getElementById('efficiencyAvg').textContent = efficiencyAvg.toFixed(1);
const evaporationAvg = data.reduce((sum, d) => sum + d.Evaporation, 0) / data.length;
document.getElementById('evaporationAvg').textContent = evaporationAvg.toFixed(1);

new Chart(evaporationChart, {
    type: 'bar',
    data: {
        labels: data.map(row => row.FinishedTimeString),
        datasets: [
            {
                label: 'Evaporation',
                data: data.map(row => row.Evaporation),
                customText: data.map(row => row.RecipeName)
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