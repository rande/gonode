/**
 * Pack revisions by month.
 *
 * @param {Array} revisions
 * @returns {Array}
 */
export default function nodeRevisionsByMonth(revisions) {
    const months = [];

    let currentMonth;
    let currentMonthId;

    revisions.forEach(revision => {
        const revDate    = new Date(revision.updated_at);
        const revMonthId = `${revDate.getMonth()}.${revDate.getFullYear()}`;

        if (revMonthId !== currentMonthId) {
            currentMonth = {
                month: revDate.getMonth(),
                year:  revDate.getFullYear(),
                items: []
            };
            months.push(currentMonth);
            currentMonthId = revMonthId;
        }

        currentMonth.items.push(revision);
    });

    return months;
}
