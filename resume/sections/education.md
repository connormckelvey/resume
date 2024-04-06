### Education {pStyle="SectionHeading"}

<%
    education.map(item => template("education-item.md", item))
        .join("\n")
%>
