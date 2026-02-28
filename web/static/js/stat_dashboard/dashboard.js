function createBarChart(canvasId, label, valueKey, data, averageElementId) {
    // Filter and sort data
    const filtered = data.filter(row => row[valueKey] !== null && !isNaN(row[valueKey]))
    // Calculate average
    const average = filtered.reduce((sum, row) => sum + row[valueKey], 0) / filtered.length;
    document.getElementById(averageElementId).textContent = average.toFixed(1);

    // Create chart
    new Chart(document.getElementById(canvasId), {
        type: 'bar',
        data: {
            labels: filtered.map(row => row.FinishedTimeString),
            datasets: [
                {
                    label: label,
                    data: filtered.map(row => row[valueKey]),
                    customText: filtered.map(row => row.RecipeName)
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
}

const rawData = window.Stats;
createBarChart('efficiencyChart', 'Efficiency', 'Efficiency', rawData, 'efficiencyAvg');
createBarChart('evaporationChart', 'Evaporation', 'Evaporation', rawData, 'evaporationAvg');

