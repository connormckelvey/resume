### Skills {pStyle="SectionHeading" class="SectionHeading"}

<%
    const categories = skills
        .reduce((categories, skill) => {
            const category = categories[skill.category] ?? []
            return {
                ...categories,
                [skill.category]: [...category, skill]
            }
        }, {})

    Object.entries(categories).map(([category, skills]) => {
        return skills
            .map((skill) => {
                return template("../util/definition-list.md", {
                    term: `${skill.name}`,
                    definition: skill.keywords.join(", ")
                })
            })
            .join("\n")
    })
    .join("\n")
%>
