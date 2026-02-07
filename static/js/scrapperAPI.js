document.getElementById("scrape-form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const urlsText = e.target[0].value;

	// Split by newlines and filter out empty lines, then joins them with a comma
	const urls = urlsText
		.split("\n")
		.map((url) => url.trim())
		.filter((url) => url.length > 0)
		.join(",");

	console.log(`Processing ${urls.split(",").length} URL(s)`);

	// Process each URL
	try {
		const res = await fetch("/api/v1/scrape", {
			method: "POST",
			headers: { "content-type": "text/plain" },
			body: urls,
		});
		const data = await res.text();
		console.log(`Scraped ${urls}:`, data);
	} catch (error) {
		console.error(`Error scraping ${urls}:`, error);
	}

	closeModal();
});

// Modal functionality
function openModal() {
	document.getElementById("scrapeModal").style.display = "block";
}

function closeModal() {
	document.getElementById("scrapeModal").style.display = "none";
	// Clear the form
	document.getElementById("scrape-form").reset();
	const countDisplay = document.getElementById("url-count");
	if (countDisplay) {
		countDisplay.textContent = "0 URL(s) entered";
	}
}

// Close modal when clicking outside
window.onclick = function (event) {
	const modal = document.getElementById("scrapeModal");
	if (event.target == modal) {
		closeModal();
	}
};
