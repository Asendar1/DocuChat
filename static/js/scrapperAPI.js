document.getElementById("scrape-form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const url = e.target[0].value;
	const res = await fetch("/api/v1/scrape", {
		method: "POST",
		headers: {"content-type": "text/plain"},
		body: url
	})
	const data = await res.text();
	console.log(data);
})
