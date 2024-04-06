#### <% name %>
*<% institution %> — <% template("../util/location.md", location) %> [ .] <% template("../util/dates.md", { dates }) %>*

---

<% 
    highlights
        .map(h => `- ${h}`)
        .join("\n")
%>

---
