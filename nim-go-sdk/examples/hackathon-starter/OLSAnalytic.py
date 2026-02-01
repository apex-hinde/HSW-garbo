import torch
import torch.nn as nn
import torch.nn.functional as F
from torch import linalg as LA

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

from sklearn import linear_model
from sklearn.preprocessing import StandardScaler

from tqdm import tqdm
import os
import sys
import warnings
warnings.filterwarnings('ignore')

from sklearn.metrics import (
    accuracy_score,
    mean_squared_error,
    mean_absolute_error,
    r2_score,
    root_mean_squared_error
)

import random
SEED = 42
torch.manual_seed(SEED)
torch.cuda.manual_seed(SEED)
torch.cuda.manual_seed_all(SEED)
np.random.seed(SEED)
random.seed(SEED)
torch.backends.cudnn.deterministic = True
torch.backends.cudnn.benchmark = False

class OLSAnalytical(nn.Module):
    """
    OLS using the analytical solution: β (X^T X)^(-1) X^T y
    This computes the optimal weights directly without gradient descent.
    """
    def __init__(self, in_features, out_features):
        super(OLSAnalytical, self).__init__()

        self.weight = nn.Parameter(torch.zeros(out_features, in_features))
        self.bias = nn.Parameter(torch.zeros(out_features))

    def fit(self, X, y):
        """
        Fit the model using the analytical OLS solution.
        X: (n_samples, in_features)
        y: (n_samples, out_features)
        """
        # Add bias column to X
        X_with_bias = torch.cat([torch.ones(X.shape[0], 1), X], dim=1)

        # Analytical solution: β = (X^T X)^(-1) X^T y
        XtX = X_with_bias.t() @ X_with_bias
        Xty = X_with_bias.t() @ y

        # Beta calculated with least squares solution
        beta = LA.solve(XtX, Xty)

        # Set parameters (no gradient needed)
        with torch.no_grad():
            self.bias.copy_(beta[0])
            self.weight.copy_(beta[1:].t())

    def forward(self, x):

        return x @ self.weight.t() + self.bias

if __name__ == "__main__":
    device = torch.device("cpu")

    model = OLSAnalytical(in_features=1, out_features=1).to(device) 

    torch.onnx.export(
        model,
        torch.randn(1, 1).to(device),
        "ols_model.onnx",
        input_names=["Dates"],
        output_names=["Income"],
        dynamic_axes={
            "input": {0: "batch_size"},
            "output": {0: "batch_size"}
        },
        opset_version=13
    )