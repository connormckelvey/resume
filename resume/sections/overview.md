# <% name %>
## <% title %>

<% template("../util/location.md", location) %>

<% 
    [{ text: email, href:`mailto:${email}` }, ...links]
        .map(link => `[${link.text}](${link.href})`)
        .join(" | ")
%>

<% summary %>
