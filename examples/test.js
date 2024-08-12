import beanstalkd from 'k6/x/beanstalkd';
import { check, sleep } from 'k6';

function cleanupTube(client, tubeName) {
    console.log(`Cleaning up tube: ${tubeName}`);
    client.watch(tubeName);
    while (true) {
        try {
            const [jobId, _] = client.reserve(0);
            client.delete(jobId);
            console.log(`Deleted job ${jobId} from ${tubeName}`);
        } catch (error) {
            if (error.message.includes("timeout")) {
                break;
            } else {
                console.error(`Error cleaning up ${tubeName}:`, error);
                break;
            }
        }
    }
    if (tubeName !== "default") {
        client.ignore(tubeName);
    }
}

export default function () {
    console.log('Starting Beanstalkd test...');
    
    const client = beanstalkd.newClient('localhost:11300');
    console.log('Client created');
    
    check(client, {
        'client is created': (c) => c !== null,
    });
    
    try {
        cleanupTube(client, "default");
        cleanupTube(client, "test-tube");
        
        const jobContent = 'Hello, Beanstalkd!';
        const jobId = client.put(jobContent, 1, 0, 60);
        console.log(`Put job with ID: ${jobId} in default tube`);
        
        const [reservedId, jobBody] = client.reserve(5);
        console.log(`Reserved job ${reservedId} from default tube: ${jobBody}`);
        
        check(jobBody, {
            'reserved job content matches': (body) => body === jobContent,
        });
        
        client.delete(reservedId);
        console.log('Job deleted from default tube');
        
        console.log('Testing tube management...');

        let tubes = client.listTubes();
        console.log('Initial tubes:', tubes);
        check(tubes, {
            'default tube exists': (t) => t.includes('default'),
        });

        const newTube = 'test-tube';
        console.log(`Using new tube: ${newTube}`);
        client.use(newTube);

        const newTubeJobContent = 'Job in new tube';
        const newTubeJobId = client.put(newTubeJobContent, 1, 0, 60);
        console.log(`Put job in ${newTube} with ID: ${newTubeJobId}`);

        tubes = client.listTubes();
        console.log('Tubes after creating new tube:', tubes);
        check(tubes, {
            'new tube was created': (t) => t.includes(newTube),
        });
        
        console.log(`Watching tube: ${newTube}`);
        client.watch(newTube);

        console.log('Attempting to reserve from the new tube');
        const [newTubeReservedId, newTubeJobBody] = client.reserve(5);
        console.log(`Reserved job from ${newTube}: ${newTubeReservedId}, content: ${newTubeJobBody}`);
        
        check(newTubeReservedId, {
            'reserved job from new tube': (id) => id === newTubeJobId,
        });
        
        check(newTubeJobBody, {
            'reserved job content from new tube matches': (body) => body === newTubeJobContent,
        });

        if (newTubeReservedId) {
            client.delete(newTubeReservedId);
            console.log(`Deleted job ${newTubeReservedId} from ${newTube}`);
        }

        console.log(`Ignoring tube: ${newTube}`);
        client.ignore(newTube);

        console.log('Testing statistics...');
        const serverStats = client.stats();
        console.log('Server stats:', JSON.stringify(serverStats, null, 2));
        
        const tubeStats = client.statsTube(newTube);
        console.log(`Stats for tube ${newTube}:`, JSON.stringify(tubeStats, null, 2));
        
    } catch (error) {
        console.error('Error occurred:', JSON.stringify(error, null, 2));
    } finally {
        if (client && client.close) {
            client.close();
        }
    }
    
    sleep(1);
}