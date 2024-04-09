#### <% name %>
*<% company  %> — <% template("../util/location.md", location) %><% this.remote ? " / Remote" : "" %> [ .] <% template("../util/dates.md", { dates }) %>*

<% 
    highlights
        .map(h => `- ${h}`)
        .join("\n")
%>
