interface ProblemDetail {
    id: number;
    title: string;
    statement: string;
    answer: string;
    author: string;
    reference: string;
    referenceURL: string;
    createdAt: Date;
}

export default ProblemDetail
