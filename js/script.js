const getData = async () => {
    try {
        const target = `data.csv`; //file
        const res = await fetch(target, {
            method: 'get',
            headers: {
                'content-type': 'text/csv;charset=UTF-8',
                //'Authorization': //in case you need authorisation
            }
        });
        if (res.status === 200) {
            const data = await res.text();
            console.log(data);
            g = new Dygraph(
                document.getElementById("graphdiv"),
                data
            );
        } else {
            console.log(`Error code ${res.status}`);
        }
    } catch (err) {
        console.log(err)
    }
}

// Run
getData();