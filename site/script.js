// No backend wired up — this only shows a local confirmation so the form
// isn't a dead end during design review. Replace with a real submit
// handler (Formspree, a serverless function, etc.) before using this
// site for actual signups. See the form-note text in index.html.
document.getElementById("waitlist-form")?.addEventListener("submit", function (e) {
  e.preventDefault();
  const form = e.currentTarget;
  const name = form.querySelector("#name").value.trim();

  const success = document.createElement("div");
  success.className = "form-success";
  success.textContent = name
    ? `Thanks, ${name} — this is a design placeholder, no request was actually sent. Wire this form to a real endpoint first.`
    : "Thanks — this is a design placeholder, no request was actually sent. Wire this form to a real endpoint first.";

  form.replaceWith(success);
});
