import TargetList from '../components/TargetList';
import { getDBConnection } from '../lib/db.js'

export default TargetList;

export async function getServerSideProps({ query }) {
    try {
        let { page } = query;

        let searchString = "";

        page = "0" + page

        let xor = 0x40
        for (let i = 0; i < page.length; i += 2) {
            const pageNumber = page.slice(i, i + 2)

            console.log(pageNumber)

            const char = parseInt(pageNumber, 16);
            const asciiCode = xor ^ char;
            const prefix = String.fromCharCode(asciiCode);
            console.log(prefix)

            searchString = searchString + prefix;


        }
        const searchQuery = `SELECT * FROM targets WHERE name LIKE "${searchString}%";`
        console.log(searchQuery)
        const db = getDBConnection();
        const [targets] = await db.query(
            searchQuery
        );

        return {
            props: {
                targets,
                pageId: page,
                searchString: searchString,
            },
        };
    } catch {
        return {
            props: {
                targets: [],
                pageId: 'failed',
                searchString: 'failed',
            },
        };
    }
}
